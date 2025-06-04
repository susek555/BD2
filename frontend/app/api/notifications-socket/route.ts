import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { NOTIFICATIONS_SOCKET_URL } from "@/app/lib/constants";


export async function GET(){
  const response = await fetchWithRefresh(`${NOTIFICATIONS_SOCKET_URL}`, {
    method: "GET",
  });
  if (!response.ok) {
    console.error("Failed to connect to notifications socket:", response.statusText);
    return Response.json({ error: "Failed to connect to notifications socket" }, { status: 511 });
  }
  return Response.json({});
}
