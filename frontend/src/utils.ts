export function datePrettyPrint(isoDate: string) {
  const date = new Date(isoDate);

  const day = String(date.getUTCDate()).padStart(2, "0");
  const month = String(date.getUTCMonth() + 1).padStart(2, "0");
  const year = date.getUTCFullYear();

  let hours = date.getUTCHours();
  const minutes = String(date.getUTCMinutes()).padStart(2, "0");
  const ampm = hours >= 12 ? "PM" : "AM";
  hours = hours % 12 || 12;
  const hourStr = String(hours).padStart(2, "0");

  const formatted = `${month}/${day}/${year}, ${hourStr}:${minutes} ${ampm}`;

  return formatted;
}

export function simpleDatePrettyPrint(isoDate: string) {
  const date = new Date(isoDate);

  const day = String(date.getUTCDate()).padStart(2, "0");
  const month = String(date.getUTCMonth() + 1).padStart(2, "0");
  const year = date.getUTCFullYear();

  const formatted = `${month}/${day}/${year}`;

  return formatted;
}

export function getTimeAgo(dateString: string): string {
  const now = new Date();
  const past = new Date(dateString);
  const diffMs = now.getTime() - past.getTime();

  const msPerHour = 1000 * 60 * 60;
  const msPerDay = msPerHour * 24;
  const msPerMonth = msPerDay * 30.44; // average
  const msPerYear = msPerDay * 365.25; // average

  if (diffMs < msPerDay) {
    const hours = Math.floor(diffMs / msPerHour);
    return `${hours} hour${hours !== 1 ? "s" : ""} ago`;
  } else if (diffMs < msPerMonth) {
    const days = Math.floor(diffMs / msPerDay);
    return `${days} day${days !== 1 ? "s" : ""} ago`;
  } else if (diffMs < msPerYear) {
    const months = Math.floor(diffMs / msPerMonth);
    return `${months} month${months !== 1 ? "s" : ""} ago`;
  } else {
    const years = Math.floor(diffMs / msPerYear);
    return `${years} year${years !== 1 ? "s" : ""} ago`;
  }
}