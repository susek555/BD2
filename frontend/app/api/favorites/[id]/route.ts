import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { API_URL } from "@/app/lib/constants";

export async function POST(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;

  await fetchWithRefresh(`${API_URL}/favourites/like/${id}`, {
    method: "POST",
  });

    return Response.json({ success: true });
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;

  const response = await fetchWithRefresh(`${API_URL}/favourites/dislike/${id}`, {
    method: "DELETE",
  });

  if(response.status === 200) {
    return Response.json({ success: true });
  } else {
    return Response.json({ success: false })
  }
}
