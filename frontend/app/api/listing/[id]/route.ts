import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/id/${id}`,
  );

  return response.json();
}
