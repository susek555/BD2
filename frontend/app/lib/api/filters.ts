const API_URL = process.env.API_URL;

export async function getColors() : Promise<string[]> {
  const response = await fetch(`${API_URL}/car/colors`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch colors");
  }

  const data = await response.json();

  return data.colors;
}

export async function getDrives() : Promise<string[]> {
  const response = await fetch(`${API_URL}/car/drives`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch drives");
  }

  const data = await response.json();

  return data.drives;
}

export async function getFuelTypes() : Promise<string[]> {
  const response = await fetch(`${API_URL}/car/fuel-types`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch fuel types");
  }

  const data = await response.json();

return data["fuel-types"];
}

export async function getTransmissions() : Promise<string[]> {
  const response = await fetch(`${API_URL}/car/transmissions`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch transmissions");
  }

  const data = await response.json();

  return data.transmissions;
}