export async function UploadImages(images: File[], id): Promise<number> {
  console.log("Posting images:", images.length);

  const formData = new FormData();
  images.forEach((image) => {
    formData.append("images", image);
  });

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