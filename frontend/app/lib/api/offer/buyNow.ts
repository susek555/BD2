export async function buyNowAuction(id: string) {

  const response = await fetch(`/api/auctions/buy-now/${id}`, {
    method: "POST",
  });

  console.log("Response status:", response.status);

  if (response.status !== 200 && response.status !== 201) {
    const errorText = await response.text();
    throw new Error(`Failed to buy now auction: ${response.status} – ${errorText}`);
  }
}

export async function buyRegular(id: string) {

  const response = await fetch(`/api/sale-offer/buy/${id}`, {
    method: "POST",
  });

  console.log("Response status:", response.status);

  if (response.status !== 200 && response.status !== 201) {
    const errorText = await response.text();
    throw new Error(`Failed to buy now auction: ${response.status} – ${errorText}`);
  }
}