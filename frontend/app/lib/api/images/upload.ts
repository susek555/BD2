export async function UploadImages(images: File[], id: number): Promise<number> {
  console.log("Posting images:", images.length);

  const formData = new FormData();
  images.forEach((image) => {
    formData.append("images", image);
  });

  console.log("Uploading images for ID:", id);

  const response = await fetch(`/api/images/upload/${id}`, {
  method: "PATCH",
  body: formData,
});

  if (response.status === 201) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to upload images: ${response.status} – ${errorText}`);
  }
}