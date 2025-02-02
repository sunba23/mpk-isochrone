import { StopsDetailsResponse, parseLocation } from "../types/stops";
import { TravelDataResponse } from "../types/traveldata";

export const fetchStops = async (): Promise<StopData[]> => {
  try {
    const response = await fetch("http://localhost:8080/stops/details", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
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

export const fetchTravelData = async (
  stopId: number,
): Promise<TravelDataResponse> => {
  try {
    const response = await fetch(
      `http://localhost:8080/traveldata?stop_id=${stopId}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      },
    );
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching travel data:", error);
    throw error;
  }
};
