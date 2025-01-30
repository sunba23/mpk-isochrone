package scripts

var (
	st1 = `
		CREATE TABLE IF NOT EXISTS route_edges (
			id SERIAL PRIMARY KEY,
			source TEXT,
			target TEXT,
			travel_time INT,
			CONSTRAINT fk_source FOREIGN KEY (source) REFERENCES stops(stop_id),
			CONSTRAINT fk_target FOREIGN KEY (target) REFERENCES stops(stop_id)
		);
  `
	st2 = `
  WITH consecutive_stops AS (
    SELECT 
        t1.trip_id,
        t1.stop_id as source_stop,
        t2.stop_id as target_stop,
        (t1.arrival_time + '00:00:00'::time)::time as source_arrival,
        (t2.arrival_time + '00:00:00'::time)::time as target_arrival
    FROM stop_times t1
    JOIN stop_times t2 
        ON t1.trip_id = t2.trip_id 
        AND t2.stop_sequence = t1.stop_sequence + 1
  )
  INSERT INTO route_edges (source, target, travel_time)
  SELECT 
      source_stop,
      target_stop,
      EXTRACT(EPOCH FROM (target_arrival - source_arrival))::integer as travel_time
  FROM consecutive_stops
  GROUP BY source_stop, target_stop, 
           EXTRACT(EPOCH FROM (target_arrival - source_arrival))::integer
  ON CONFLICT DO NOTHING;
  `
	st3 = `
		CREATE TABLE IF NOT EXISTS precomputed_travel_times AS
		SELECT 
			start.stop_id AS from_stop_id,
			result.node AS to_stop_id,
			result.cost AS travel_time
		FROM stops start,
		LATERAL (
			SELECT * FROM pgr_dijkstra(
				'SELECT id, source, target, travel_time AS cost FROM route_edges',
				start.stop_id,
				ARRAY(SELECT stop_id FROM stops WHERE stop_id != start.stop_id),
				directed := true
			)
		) AS result;
	`
)
