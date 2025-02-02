import React, { useState, useEffect, useRef } from "react";
import { MapContainer, TileLayer, CircleMarker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { fetchStops, fetchTravelData } from "../services/apiService";
import { StopData } from "../types/stops";
import { getMarkerColor } from "../utils/color";
import MapLegend from "./MapLegend";
import { FaPlay, FaFlagCheckered, FaBus } from "react-icons/fa";
import { FaTrainTram } from "react-icons/fa6";

const IsochronicMap: React.FC = () => {
  const [stops, setStops] = useState<StopData[]>([]);
  const [markerColors, setMarkerColors] = useState<Record<string, string>>({});
  const [travelTimes, setTravelTimes] = useState<Record<string, number>>({});
  const [startStopId, setStartStopId] = useState<string | null>(null);
  const markersRef = useRef<Record<string, L.CircleMarker>>({});
  const [endStopId, setEndStopId] = useState<string | null>(null);
  const [pathStops, setPathStops] = useState<string[]>([]);
  const [paths, setPaths] = useState<Record<string, string[]>>({});

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

      const paths: Record<string, string[]> = {};

      Object.entries(travelData.stop_id_travel_data_map).forEach(
        ([destStopId, data]) => {
          times[destStopId] = data.travel_time;

          if (data.path && data.path.length > 0) {
            paths[destStopId] = data.path;
          }
        },
      );

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
      setEndStopId(null);
      setPaths(paths);
      setPathStops([]);
    } catch (error) {
      console.error("Error generating isochrone:", error);
    }
  };

  useEffect(() => {
    if (endStopId && paths[endStopId]) {
      setPathStops(paths[endStopId]);
    }
  }, [endStopId, paths]);

  useEffect(() => {
    Object.entries(markersRef.current).forEach(([stopId, marker]) => {
      if (!endStopId) {
        // Isochrone mode
        const isStart = stopId === startStopId;
        marker.setStyle({
          fillColor: isStart ? "#FF00FF" : markerColors[stopId] || "#CCCCCC",
          radius: isStart ? 7 : 5,
          color: "#000000",
          fillOpacity: isStart ? 0.95 : 0.85,
          weight: isStart ? 3 : 1,
          zIndexOffset: isStart ? 1000 : 0,
        });
      } else {
        // Finish mode
        const isStart = stopId === startStopId;
        const isEnd = stopId === endStopId;
        const isInPath = pathStops.includes(Number(stopId));

        // uninmportant stops
        let style = {
          fillColor: markerColors[stopId] || "#CCCCCC",
          radius: 4,
          color: "rgba(0,0,0,0.3)",
          weight: 0.5,
          fillOpacity: 0.3,
          zIndexOffset: 0,
        };

        // important stops
        if (isStart) {
          style = {
            fillColor: "#FF00FF",
            radius: 9,
            color: "#000000",
            weight: 3,
            fillOpacity: 0.95,
            zIndexOffset: 3000,
          };
        } else if (isEnd) {
          style = {
            fillColor: "#00FF00",
            radius: 8,
            color: "#000000",
            weight: 3,
            fillOpacity: 0.9,
            zIndexOffset: 2000,
          };
        } else if (isInPath) {
          style = {
            fillColor: markerColors[stopId] || "#CCCCCC",
            radius: 7,
            color: "#000000",
            weight: 2,
            fillOpacity: 0.85,
            zIndexOffset: 1000,
          };
        }

        marker.setStyle(style);

        const element = marker.getElement();
        if (element) {
          if (isStart || isEnd || isInPath) {
            element.classList.add("important-marker");
          } else {
            element.classList.remove("important-marker");
          }
        }
      }
    });
  }, [markerColors, startStopId, endStopId, pathStops]);

  return (
    <div className="main">
      <div
        style={{
          width: "80vw",
          maxWidth: "1200px",
          margin: "0 auto",
        }}
      >
        <h1
          className="map-title text-center"
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: "8px",
          }}
        >
          Wrocław public transport isochrone/travel visualizer
          <b>
            <FaBus style={{ verticalAlign: "middle", margin: "0 4px" }} />
            <FaTrainTram style={{ verticalAlign: "middle", margin: "0 4px" }} />
          </b>
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
            preferCanvas={true}
            zoomControl={true}
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
                    aria-label="Start"
                  >
                    <FaPlay style={{ marginRight: 5 }} />
                  </button>
                  {startStopId && (
                    <button
                      onClick={() => setEndStopId(stop.id)}
                      aria-label="Finish"
                    >
                      <FaFlagCheckered style={{ marginLeft: 5 }} />
                    </button>
                  )}
                </Popup>
              </CircleMarker>
            ))}
          </MapContainer>
        </div>
        <p className="mt-4 text-gray-700 leading-relaxed map-desc">
          This map visualizes travel times for Wrocław public transport. To
          start, <b>click on a stop</b> and then{" "}
          <b>
            <FaPlay style={{ verticalAlign: "middle", margin: "0 4px" }} />
          </b>
          . The travel times are <b>best case scenario</b> (e.g. instant
          transfer). To see the stops needed to get somewhere,{" "}
          <b>click on another stop</b> and then{" "}
          <b>
            <FaFlagCheckered
              style={{ verticalAlign: "middle", margin: "0 4px" }}
            />
          </b>
          . Like it? Star the project on{" "}
          <a href="https://github.com/sunba23/mpk-isochrone">GitHub</a>!
        </p>
      </div>
    </div>
  );
};

export default IsochronicMap;
