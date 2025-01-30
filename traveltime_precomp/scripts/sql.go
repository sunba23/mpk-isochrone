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
      t1.departure_time::interval as source_departure,
      t2.arrival_time::interval as target_arrival
    FROM stop_times t1
    JOIN stop_times t2 
      ON t1.trip_id = t2.trip_id 
      AND t2.stop_sequence = t1.stop_sequence + 1
  )
  INSERT INTO route_edges (source, target, travel_time)
  SELECT 
    source_stop,
    target_stop,
    CASE 
      WHEN target_arrival < source_departure 
      THEN EXTRACT(EPOCH FROM (target_arrival + interval '24 hours' - source_departure))::integer
      ELSE EXTRACT(EPOCH FROM (target_arrival - source_departure))::integer
    END as travel_time
  FROM consecutive_stops
  GROUP BY source_stop, target_stop, 
           CASE 
             WHEN target_arrival < source_departure 
             THEN EXTRACT(EPOCH FROM (target_arrival + interval '24 hours' - source_departure))::integer
             ELSE EXTRACT(EPOCH FROM (target_arrival - source_departure))::integer
           END
  ON CONFLICT DO NOTHING;
  `
  st3 = `
  INSERT INTO route_edges (source, target, travel_time)
  SELECT sg1.stop_id, sg2.stop_id, 0
  FROM stop_groups sg1
  JOIN stop_groups sg2 ON sg1.group_id = sg2.group_id
  WHERE sg1.stop_id <> sg2.stop_id;

  INSERT INTO route_edges (source, target, travel_time)
  SELECT sg2.stop_id, sg1.stop_id, 0
  FROM stop_groups sg1
  JOIN stop_groups sg2 ON sg1.group_id = sg2.group_id
  WHERE sg1.stop_id <> sg2.stop_id;
  `
	st4 = `
  CREATE TABLE IF NOT EXISTS precomputed_travel_times AS
  WITH paths AS (
    SELECT 
      start.stop_id AS from_stop_id,
      result.end_vid AS to_stop_id,
      MAX(result.agg_cost) AS travel_time,
      array_agg(result.node ORDER BY result.path_seq) AS route_stops
    FROM stops start,
    LATERAL (
      SELECT 
        path_seq,
        node,
        end_vid,
        agg_cost,
        edge
      FROM pgr_dijkstra(
        'SELECT id,
                source::integer as source,
                target::integer as target,
                travel_time AS cost
         FROM route_edges',
        start.stop_id::int,
        ARRAY(SELECT stop_id FROM stops WHERE stop_id != start.stop_id),
        directed := true
      )
      WHERE edge != -1
    ) AS result
    GROUP BY start.stop_id, result.end_vid
  )
  SELECT * FROM paths WHERE travel_time > 0;
  `
)
