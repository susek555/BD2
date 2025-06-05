export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  const response = await fetch(
    `${process.env.API_URL}/review/average-rating/${id}`,
  );
  const data = await response.json();
  return Response.json(data);
}
