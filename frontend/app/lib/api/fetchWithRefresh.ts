import { authConfig } from '@/app/lib/authConfig';
import { getServerSession } from 'next-auth';
import { getSession, signIn } from 'next-auth/react';

export async function fetchWithRefresh(
  input: RequestInfo,
  init: RequestInit = {},
): Promise<Response> {
  const session = await getServerSession(authConfig);
  const accessToken = session?.user?.accessToken;

  let res = await fetch(input, {
    ...init,
    headers: {
      ...(init.headers ?? {}),
      Authorization: `Bearer ${accessToken}`,
    },
  });

  // TODO fix handling when token expired
  if (res.status === 401 || res.status === 403) {
    await signIn('credentials', { redirect: false });

    const newSession = await getSession();
    const newToken = newSession?.user?.accessToken;

    res = await fetch(input, {
      ...init,
      headers: {
        ...(init.headers ?? {}),
        Authorization: `Bearer ${newToken}`,
      },
    });
  }

  return res;
}
