#!/usr/bin/env python3
import logging
import os
import shutil
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path

import requests
from rich.console import Console
from rich.logging import RichHandler

GTFS_BINARY_URL = "https://github.com/public-transport/gtfs-via-postgres/releases/download/4.10.2/gtfs-via-postgres-linux-x64"

@dataclass
class Config:
    gtfs_url: str
    output_dir: Path
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
            handlers=[RichHandler(console=self.console)]
        )
        self.logger = logging.getLogger("gtfs_importer")

    def ensure_binary(self):
        if not self.config.binary_path.exists():
            self.logger.info("Downloading gtfs-via-postgres...")
            response = requests.get(GTFS_BINARY_URL)
            response.raise_for_status()
            
            self.config.binary_path.write_bytes(response.content)
            self.config.binary_path.chmod(0o755)

    def download_gtfs(self) -> Path:
        self.logger.info("Downloading GTFS data...")
        response = requests.get(self.config.gtfs_url, stream=True)
        response.raise_for_status()
        
        zip_path = self.config.output_dir / "gtfs.zip"
        with open(zip_path, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
        return zip_path

    def extract_gtfs(self, zip_path: Path) -> Path:
        extract_dir = self.config.output_dir / "gtfs"
        shutil.unpack_archive(zip_path, extract_dir)
        return extract_dir

    def import_gtfs(self, gtfs_path: Path):
        cmd = [
            str(self.config.binary_path),
            "--database", self.config.db_url,
            str(gtfs_path)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            raise RuntimeError(f"Import failed: {result.stderr}")

    def run(self):
        try:
            self.ensure_binary()
            zip_path = self.download_gtfs()
            gtfs_path = self.extract_gtfs(zip_path)
            self.import_gtfs(gtfs_path)
            self.logger.info("GTFS import completed successfully!")
        except Exception as e:
            self.logger.error(f"GTFS import failed: {e}")
            sys.exit(1)

if __name__ == "__main__":
    config = Config(
        gtfs_url=os.environ["GTFS_URL"],
        output_dir=Path(os.getenv("OUTPUT_DIR", "/tmp/gtfs")),
        binary_path=Path(os.getenv("BINARY_PATH", "/usr/local/bin/gtfs-via-postgres")),
        db_url=os.environ["DATABASE_URL"]
    )

    if not config.gtfs_url or not config.db_url:
        logging.error("GTFS_URL and DATABASE_URL environment variables are required")
        sys.exit(1)

    importer = GTFSImporter(config)
    importer.run()
