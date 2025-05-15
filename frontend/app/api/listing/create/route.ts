const API_URL = process.env.API_URL;

export async function POST(request: Request) {
  console.log('Create request');
  // TODO implement API call to add listing
  return Response.json({ success: true });
}
