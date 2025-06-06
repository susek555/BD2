export async function uploadImages(images: File[], id: number): Promise<number> {
  console.log("Posting images:", images.length);

  const formData = new FormData();
  images.forEach((image) => {
    formData.append("images", image);
  });

  // Log formData contents for debugging
  for (const pair of formData.entries()) {
    console.log(pair[0], pair[1]);
  }
  console.log("images", images);
  console.log("Uploading images for ID:", id);

  const response = await fetch(`/api/images/upload/${id}`, {
    method: "PATCH",
    body: formData,
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


