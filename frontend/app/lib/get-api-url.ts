export function getApiUrl(relative: string): string {
  const isServer = typeof window === 'undefined';

  if (isServer) {
    const baseUrl = process.env.NEXTAUTH_URL
      ? `${process.env.NEXTAUTH_URL}`
      : 'http://localhost:3000';

    return `${baseUrl}${relative}`;
  } else {
    return relative;
  }
}
