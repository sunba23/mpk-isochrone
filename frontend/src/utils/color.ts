export const getMarkerColor = (travelTime: number, minTime: number, maxTime: number): string => {
  const normalizedTime = (travelTime - minTime) / (maxTime - minTime);

  if (normalizedTime <= 0.14) {
    return '#FF0000';
  } else if (normalizedTime <= 0.28) {
    return '#FFA600'
  } else if (normalizedTime <= 0.42) {
    return '#F9FF00';
  } else if (normalizedTime <= 0.56) {
    return '#28FF00';
  } else if (normalizedTime <= 0.70) {
    return '#00FFF8';
  } else if (normalizedTime <= 0.86) {
    return '#008CFF';
  } else {
    return '#11007F';
  }
};
