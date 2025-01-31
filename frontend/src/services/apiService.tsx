import { StopsDetailsResponse, parseLocation } from "../types/stops";

export const fetchStops = async (): Promise<StopData[]> => {
  try {
    const response = await fetch("http://localhost:8080/stops/details", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    console.log("Response status:", response.status);
    const rawText = await response.text();
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    const data: StopsDetailsResponse = JSON.parse(rawText);

    return Object.values(data.stop_details_map).map((stop) => ({
      ...stop,
      stop_location: parseLocation(stop.stop_location),
    }));
  } catch (error) {
    console.error("Full fetch error: ", error);
    throw error;
  }
};
