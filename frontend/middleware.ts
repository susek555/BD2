import { getSession } from "next-auth/react";
import { NextRequest } from "next/server";
import { LOGIN, PROTECTED_ROUTES, PUBLIC_ROUTES } from "./app/lib/routes";

export async function middleware(reqest: NextRequest) {
  const { nextUrl } = reqest;
  const session = await getSession();
  const isAuthenticated = !!session?.user;

  console.log("Is authenticated:", isAuthenticated);

  const isPublicRoute = ((PUBLIC_ROUTES.find(route => nextUrl.pathname.startsWith(route)))
    && !PROTECTED_ROUTES.find(route => nextUrl.pathname.includes(route)));

  console.log("Public route:", isPublicRoute);

  if (!isAuthenticated && !isPublicRoute) {
    console.log("Redirecting to /login");
    return Response.redirect(new URL(LOGIN, nextUrl));
  }

  if (session && nextUrl.pathname === '/login') {
    return Response.redirect(new URL('/', nextUrl))
  }
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico, sitemap.xml, robots.txt (metadata files)
     */
    '/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
  ],
}
