import { API_URL } from '@/app/lib/constants';
import { ExtendedJWT } from '@/types/next-auth';
import camelcaseKeys from 'camelcase-keys';
import NextAuth, { User } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

const ACCESS_TOKEN_LIFETIME = 2 * 60 * 60 * 1000; // 2 hours
// const ACCESS_TOKEN_LIFETIME = 30 * 1000; // 30 seconds

async function updateAccessToken(refreshToken: string): Promise<string> {
  try {
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      throw new Error(`Refresh token error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();

    if (!data.access_token) {
      throw new Error('No access token returned from refresh endpoint');
    }

    return data.access_token;
  } catch (error) {
    console.error('Error refreshing token:', error);
    throw error;
  }
}

// TODO: add error handling
const handler = NextAuth({
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        login: { label: "Login", type: "text" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        // Add logic here to look up the user from the credentials supplied

        const res = await fetch(`${API_URL}/auth/login`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            login: credentials?.login,
            password: credentials?.password,
          }),
        });
        const rawUser = await res.json();

        if (!rawUser.errors) {
          const user: User = camelcaseKeys(rawUser, { deep: true });
          // Any object returned will be saved in `user` property of the JWT
          return user;
        } else {
          // If you return null then an error will be displayed advising the user to check their details.
          return null;

          // You can also Reject this callback with an Error thus the user will be sent to the error page with the error message as a query parameter
        }
      }
    })
  ],
  callbacks: {
    async jwt({ token, user }) {
      console.log("JWT callback");
      const jwtToken = token as ExtendedJWT

      // If this is the first sign in, add user data to token
      if (user) {
        console.log("JWT callback initial login")
        Object.assign(jwtToken, user);
        jwtToken.accessTokenExpires = Date.now() + ACCESS_TOKEN_LIFETIME;
      }

      // Return previous token if the access token has not expired yet
      if (jwtToken.accessTokenExpires && Date.now() < jwtToken.accessTokenExpires) {
        console.log("JWT callback still valid")
        return jwtToken;
      }

      // Access token has expired, refresh it
      try {
        console.log("Token expired, refreshing...");
        const newAccessToken = await updateAccessToken(jwtToken.refreshToken);
        jwtToken.accessToken = newAccessToken;
        jwtToken.accessTokenExpires = Date.now() + ACCESS_TOKEN_LIFETIME;
        console.log(`${Math.round((jwtToken.accessTokenExpires - Date.now()) / (1000 * 60))} minutes`)
        return { ...jwtToken };
      } catch (error) {
        console.error("Error refreshing token:", error);
        return { ...token, error: "RefreshTokenError" };
      }
    },

    async session({ session, token }) {
      session.user = token as ExtendedJWT;
      return session;
    },
  },
  pages: {
    signIn: '/login',
  },
  session: {
    strategy: "jwt",
  }
})

export { handler as GET, handler as POST };
