import React, { useState, useEffect, useRef } from "react";
import { MapContainer, TileLayer, CircleMarker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { fetchStops, fetchTravelData } from "../services/apiService";
import { StopData } from "../types/stops";
import { getMarkerColor } from "../utils/color";
import MapLegend from "./MapLegend";

const IsochronicMap: React.FC = () => {
  const [stops, setStops] = useState<StopData[]>([]);
  const [markerColors, setMarkerColors] = useState<Record<string, string>>({});
  const [travelTimes, setTravelTimes] = useState<Record<string, number>>({});
  const [startStopId, setStartStopId] = useState<string | null>(null);
  const markersRef = useRef<Record<string, L.CircleMarker>>({});

  useEffect(() => {
    const loadStops = async () => {
      try {
        const fetchedStops = await fetchStops();
        setStops(fetchedStops);
      } catch (error) {
        console.error("Error loading stops:", error);
      }
    };
    loadStops();
  }, []);

  const handleGenerateIsochrone = async (stopId: number) => {
    try {
      const travelData = await fetchTravelData(stopId);

      const times: Record<string, number> = {
        [stopId.toString()]: 0,
      };

      Object.values(travelData.stop_id_travel_data_map).forEach((data) => {
        times[data.id.toString()] = data.travel_time;
      });

      const travelTimesArray = Object.values(times);
      const minTime = Math.min(...travelTimesArray);
      const maxTime = Math.max(...travelTimesArray);

      const colors: Record<string, string> = {};
      Object.entries(times).forEach(([stopId, time]) => {
        colors[stopId] = getMarkerColor(time, minTime, maxTime);
      });

      setTravelTimes(times);
      setMarkerColors(colors);
      setStartStopId(stopId.toString());
    } catch (error) {
      console.error("Error generating isochrone:", error);
    }
  };

  useEffect(() => {
    Object.entries(markersRef.current).forEach(([stopId, marker]) => {
      if (markerColors[stopId]) {
        marker.setStyle({
          fillColor: stopId === startStopId ? "#FF00FF" : markerColors[stopId],
          radius: stopId === 5,
          color: stopId === startStopId ? "#FF00FF" : "#000000",
        });
      }
    });
  }, [markerColors, startStopId]);

  return (
    <div className="main">
      <div style={{ width: "80vw" }}>
        <h1 className="text-2xl font-bold mb-4 map-title text-center">
          Wrocław public transport isochrone/travel visualizer
        </h1>
        <div
          style={{
            height: "80vh",
            width: "100%",
            borderRadius: "20px",
            overflow: "hidden",
          }}
        >
          <MapContainer
            center={[51.108, 17.028]}
            zoom={12}
            style={{ height: "100%", width: "100%" }}
          >
            <TileLayer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />
            {Object.keys(travelTimes).length > 0 && (
              <MapLegend
                minTime={Math.min(...Object.values(travelTimes)) || 0}
                maxTime={Math.max(...Object.values(travelTimes)) || 3600}
              />
            )}
            {stops.map((stop) => (
              <CircleMarker
                key={stop.id}
                center={[
                  stop.stop_location.latitude,
                  stop.stop_location.longitude,
                ]}
                radius={5}
                fillOpacity={0.8}
                color="#000000"
                fillColor={markerColors[stop.id] || "#CCCCCC"}
                weight={1}
                ref={(marker) => {
                  if (marker) {
                    markersRef.current[stop.id] = marker;
                  }
                }}
              >
                <Popup>
                  {stop.stop_name}
                  <br />
                  Code: {stop.stop_code}
                  <br />
                  {travelTimes[stop.id] && (
                    <>
                      Travel time: ~{Math.round(travelTimes[stop.id] / 60)}{" "}
                      minutes
                      <br />
                    </>
                  )}
                  <button
                    onClick={() => handleGenerateIsochrone(parseInt(stop.id))}
                  >
                    Start Here
                  </button>
                </Popup>
              </CircleMarker>
            ))}
          </MapContainer>
        </div>
        <p className="mt-4 text-gray-700 leading-relaxed map-desc">
          This map visualizes travel times for Wrocław public transport. To
          start, <b>click on a stop</b> and then <b>Start from here</b>. The
          travel times are <b>best case scenario</b> (e.g. instant transfer). To
          see the stops needed to get somewhere, <b>click on another stop</b>{" "}
          and then <b>Finish here</b>. If you like it, star the project on{" "}
          <a href="https://github.com/sunba23/mpk-isochrone">GitHub</a>!
        </p>
      </div>
    </div>
  );
};

export default IsochronicMap;
