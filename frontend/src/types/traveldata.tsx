export interface TravelData {
  id: number;
  travel_time: number;
  path: number[];
}

export interface TravelDataResponse {
  code: number;
  stop_id_travel_data_map: Record<number, TravelData>;
}
