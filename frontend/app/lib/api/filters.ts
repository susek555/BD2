
const API_URL = process.env.API_URL;

export async function getColors() {
  const response = await fetch(`${API_URL}/car/colors`, {
    method: "GET",
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch colors");
  }

  const data = await response.json();

  return data.colors;
}