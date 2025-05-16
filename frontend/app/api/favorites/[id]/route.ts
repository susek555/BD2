export async function POST(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  // TODO implement API call to add to favorites
  console.log(`Offer ${id} added to favorites`);
  return Response.json({ success: true });
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  // TODO implement API call to delete from favorites
  console.log(`Offer ${id} removed from favorites`);
  return Response.json({ success: true });
}
