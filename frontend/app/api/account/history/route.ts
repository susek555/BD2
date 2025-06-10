import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';

export async function POST(request: Request) {
  const body = request.text();
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/purchased-offers`,
  );
  return response.json();
}
