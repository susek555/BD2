export async function GET(
  request: Request,
  ) {
  const data = {
    '1': 33,
    '2': 33,
    '3': 0,
    '4': 0,
    '5': 33,
  };  return new Response(JSON.stringify(data));
}
