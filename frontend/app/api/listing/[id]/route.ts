export async function PUT(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  console.log('Update listing id:', id);
  // TODO implement API call to update listing
  return Response.json({ success: true });
}

export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  console.log('Get listing id:', id);
  // TODO implement API call to get listing
  return Response.json({ id: id });
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  console.log('Deleted listing id:', id);
  // TODO implement API call to delete listing
  return Response.json({ success: true });
}
