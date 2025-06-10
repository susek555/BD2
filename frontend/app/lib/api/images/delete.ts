export async function deleteImages(URLs: string[]): Promise<void> {
  console.log("Deleting images:", URLs.length);

  const response = await fetch(`/api/images/delete`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ URLs }),
  });

  console.log("Response status for deletion:", response.status);

  if (response.status !== 200) {
    const errorText = await response.text();
    throw new Error(`Failed to delete images: ${response.status} – ${errorText}`);
  }
}