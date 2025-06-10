import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { NextRequest } from 'next/server';

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/review/${id}`,
    {
      method: 'DELETE',
    },
  );

  return response;
}
