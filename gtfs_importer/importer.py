import logging
import os
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from zipfile import ZipFile

import requests
from rich.console import Console
from rich.logging import RichHandler

GTFS_TO_POSTGRES_BINARY_URL = "https://github.com/public-transport/gtfs-via-postgres/releases/download/4.10.2/gtfs-via-postgres-linux-x64"
GTFS_DATA_URL = "https://www.wroclaw.pl/open-data/87b09b32-f076-4475-8ec9-6020ed1f9ac0/OtwartyWroclaw_rozklad_jazdy_GTFS.zip"


@dataclass
class Config:
    gtfs_url: str
    output_dir: Path
    binary_url: str
    binary_path: Path
    db_url: str


class GTFSImporter:
    def __init__(self, config: Config):
        self.config = config
        self.console = Console()
        self.setup_logging()

    def setup_logging(self):
        logging.basicConfig(
            level=logging.INFO,
            format="%(message)s",
            handlers=[RichHandler(console=self.console)],
        )
        self.logger = logging.getLogger("gtfs_importer")

    def ensure_binary(self) -> None:
        if not self.config.binary_path.exists():
            self.logger.info("Downloading gtfs-via-postgres...")
            response = requests.get(self.config.binary_url)
            response.raise_for_status()

            self.config.binary_path.write_bytes(response.content)
            self.config.binary_path.chmod(0o755)

    def download_gtfs(self) -> Path:
        self.logger.info("Downloading GTFS data...")
        response = requests.get(self.config.gtfs_url, stream=True)
        response.raise_for_status()

        if not os.path.isdir(self.config.output_dir):
            os.mkdir(self.config.output_dir)
        zip_path = self.config.output_dir / "gtfs.zip"
        with open(zip_path, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)

        with ZipFile(zip_path) as zf:
            if zf.testzip() is not None:
                raise ValueError("Downloaded zip file is corrupted")

        return zip_path

    def extract_gtfs(self, zip_path: Path) -> Path:
        extract_dir = self.config.output_dir / "gtfs"
        with ZipFile(zip_path) as zf:
            zf.extractall(extract_dir)
        return extract_dir

    def create_extensions(self):
        # connect to db (privileged postgres acc needed) and create extensions: postgis, pgrouting
        # before running script, make sure postGIS and pgRouting are installed on system
        pass

    def import_gtfs(self, gtfs_path: Path):
        txt_files = list(gtfs_path.glob("*.txt"))
        if not txt_files:
            raise FileNotFoundError(f"No .txt files found in {gtfs_path}")

        cmd1 = [str(self.config.binary_path), "-u", *txt_files]
        cmd2 = ["sponge"]
        cmd3 = ["psql", str(self.config.db_url), "-b"]

        try:
            with (
                subprocess.Popen(cmd1, stdout=subprocess.PIPE) as process1,
                subprocess.Popen(
                    cmd2, stdin=process1.stdout, stdout=subprocess.PIPE
                ) as process2,
                subprocess.Popen(cmd3, stdin=process2.stdout) as process3,
            ):
                if process1.stdout:
                    process1.stdout.close()
                if process2.stdout:
                    process2.stdout.close()

                return_code = process3.wait()
                if return_code != 0:
                    raise subprocess.CalledProcessError(
                        return_code,
                        f"Pipeline failed: {' | '.join([' '.join(cmd1), ' '.join(cmd2), ' '.join(cmd3)])}",
                    )

        except OSError as e:
            raise OSError(f"Failed to execute pipeline: {e}")

    def run(self):
        try:
            self.ensure_binary()
            self.create_extensions()
            zip_path = self.download_gtfs()
            gtfs_path = self.extract_gtfs(zip_path)
            self.import_gtfs(gtfs_path)
            self.logger.info("GTFS import completed successfully!")
        except Exception as e:
            self.logger.error(f"GTFS import failed: {e}")
            sys.exit(1)


if __name__ == "__main__":
    config = Config(
        gtfs_url=os.environ.get("GTFS_URL", GTFS_DATA_URL),
        output_dir=Path(os.getenv("OUTPUT_DIR", "/tmp/gtfs")),
        binary_path=Path(os.getenv("BINARY_PATH", "~/gtfs-via-postgres")).expanduser(),
        binary_url=os.environ.get(
            "GTFS_TO_POSTGRES_BINARY_URL", GTFS_TO_POSTGRES_BINARY_URL
        ),
        db_url=os.environ["DATABASE_URL"],
    )

    if not config.gtfs_url or not config.db_url or not config.binary_url:
        logging.error("DATABASE_URL environment variable is required")
        sys.exit(1)

    importer = GTFSImporter(config)
    importer.run()
