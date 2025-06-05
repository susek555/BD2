export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  // TODO uncomment once fixed on backend
  // const response = await fetch(`${process.env.API_URL}/review/${id}/average`);
  // const data = await response.json();
  // return Response.json(data);
  return Response.json(1);
}
