export async function buyNow(id: string): Promise<number> {

  const response = await fetch(`/api/auctions/buy-now/${id}`, {
    method: "DELETE",
  });

  console.log("Response status:", response.status);

  if (response.status === 200 || response.status === 201) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to upload images: ${response.status} – ${errorText}`);
  }
}