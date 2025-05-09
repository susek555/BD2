import camelcaseKeys from 'camelcase-keys';
import NextAuth from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

const API_URL = process.env.API_URL;
// const ACCEST_TOKEN_LIFETIME = 2 * 60 * 60 * 1000; // 2 hours
const ACCESS_TOKEN_LIFETIME = 30 * 1000; // 30 seconds

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
      async authorize(credentials, req) {
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
        console.log(rawUser);

        if (!rawUser.errors) {
          const user = camelcaseKeys(rawUser, { deep: true });
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
      return { ...token, ...user };
    },
    async session({ session, token, user }) {
      session.user = token as any;
      return session;
    },
  },
  pages: {
    signIn: '/login',
  },
  session: {
    strategy: "jwt",
    maxAge: 30 * 24 * 60 * 60, // 30 days
  }
})

export { handler as GET, handler as POST };
