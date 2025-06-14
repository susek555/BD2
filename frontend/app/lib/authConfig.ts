import { API_URL } from '@/app/lib/constants';
import { ExtendedJWT } from '@/types/next-auth';
import camelcaseKeys from 'camelcase-keys';
import { NextAuthOptions, User } from 'next-auth';
import CredentialsProvider from 'next-auth/providers/credentials';

const ACCESS_TOKEN_LIFETIME = 30 * 60 * 1000; // 30 minutes
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
      throw new Error(
        `Refresh token error: ${response.status} ${response.statusText}`,
      );
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
export const authConfig: NextAuthOptions = {
  providers: [
    CredentialsProvider({
      name: 'Credentials',
      credentials: {
        login: { label: 'Login', type: 'text' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials) {
        // Add logic here to look up the user from the credentials supplied

        const res = await fetch(`${API_URL}/auth/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            login: credentials?.login,
            password: credentials?.password,
          }),
        });
        const rawUser = await res.json();
        console.log('Login errors:', rawUser.errors);

        if (!rawUser.errors) {
          const user: User = camelcaseKeys(rawUser, { deep: true });
          // Any object returned will be saved in `user` property of the JWT
          return user;
        } else {
          throw new Error(JSON.stringify(rawUser.errors));
          // throw new Error(rawUser.errors);
          // You can also Reject this callback with an Error thus the user will be sent to the error page with the error message as a query parameter
        }
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user, trigger, session }) {
      console.log('JWT callback');
      const jwtToken = token as ExtendedJWT;

      if (trigger === 'update' && session) {
        const isUserProfile =
          session.selector &&
          // typeof session.id === 'number' && // uncomment when main backend branch sends id on login
          typeof session.username === 'string' &&
          typeof session.email === 'string' &&
          (session.selector === 'P' || session.selector === 'C');

        const isIdPassword =
          // typeof session.id === 'number' &&
          typeof session.password === 'string' &&
          Object.keys(session).length === 2;

        if (isUserProfile || isIdPassword) {
          Object.assign(jwtToken, session);
        }
      }

      if (user) {
        console.log('JWT callback initial login');
        Object.assign(jwtToken, user);
        jwtToken.accessTokenExpires = Date.now() + ACCESS_TOKEN_LIFETIME;
      }

      if (
        jwtToken.accessTokenExpires &&
        Date.now() < jwtToken.accessTokenExpires
      ) {
        console.log('JWT callback still valid');
        return jwtToken;
      }

      try {
        console.log('Token expired, refreshing...');
        const newAccessToken = await updateAccessToken(jwtToken.refreshToken);
        jwtToken.accessToken = newAccessToken;
        jwtToken.accessTokenExpires = Date.now() + ACCESS_TOKEN_LIFETIME;
        console.log(
          `${Math.round((jwtToken.accessTokenExpires - Date.now()) / (1000 * 60))} minutes`,
        );
        return { ...jwtToken };
      } catch (error) {
        console.error('Error refreshing token:', error);
        return { ...token, error: 'RefreshTokenError' };
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
    strategy: 'jwt',
  },
};
