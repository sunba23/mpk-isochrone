import * as wkx from "wkx";

export interface StopData {
  id: string;
  stop_code: string;
  stop_name: string;
  stop_location: {
    latitude: number;
    longitude: number;
  };
}

export interface StopsDetailsResponse {
  code: number;
  stop_details_map: Record<
    string,
    Omit<StopData, "stop_location"> & { stop_location: string }
  >;
}

export function parseLocation(base64Location: string): {
  latitude: number;
  longitude: number;
} {
  const buffer = Buffer.from(base64Location, "base64");

  const geometry = wkx.Geometry.parse(buffer);

  return {
    latitude: geometry.y,
    longitude: geometry.x,
  };
}
