export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const randomRating = Math.round((Math.random() * 4 + 1) * 100) / 100;
  return new Response(JSON.stringify(randomRating));
}
