"use server";

import { cookies } from 'next/headers';

const API_URL = process.env.API_URL;

type registerResponse = {
  access_token?: string
  refresh_token?: string
  errors?: {
    username?: string[]
    email?: string[]
    nip?: string[]
  }
}

export type registerResult = {
  success: boolean
  errors?: {
    username?: string[]
    email?: string[]
    nip?: string[]
    other?: string[]
  }
}

type loginResponse = {
  access_token?: string
  refresh_token?: string
  errors?: {
    credentials?: string[]
  }
}

export type loginResult = {
  success: boolean
  errors?: {
    credentials?: string[]
    other?: string[]
  }
}

async function saveTokensToCookies(access_token: string, refresh_token: string) {
  const cookieStore = await cookies();

  cookieStore.set('access_token', access_token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    maxAge: 60 * 60 * 2, // 2 hours
    path: '/',
  });

  cookieStore.set('refresh_token', refresh_token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    maxAge: 60 * 60 * 24 * 30, // 30 days
    path: '/',
  });

  console.log("saved coookies");
}

/**
 * Sends registration data to the backend and stores tokens if successful
 * @param userData Registration data in the required format
 * @returns Response with success status and error information
 */
export async function registerUser(userData: {
  selector: 'C' | 'P';
  username: string;
  email: string;
  password: string;
  person_name?: string;
  person_surname?: string;
  company_name?: string;
  company_nip?: string;
}): Promise<registerResult> {
  console.log(JSON.stringify(userData));

  try {
    const response = await fetch(`${API_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    });

    const data: registerResponse = await response.json();

    console.log(data);

    if (!response.ok) {
      return {
        success: false,
        errors: data.errors || {},
      };
    }

    if (data.access_token && data.refresh_token) {
      saveTokensToCookies(data.access_token, data.refresh_token);
    }

    return {
      success: true,
    };
  } catch (error: any) {
    console.error('Registration error:', error);
    return {
      success: false,
      errors: {
        other: error.message || 'An unexpected error occurred during registration'
      },
    };
  }
}

/**
 * Authenticates user and stores tokens in cookies
 * @param credentials User login credentials
 * @returns Response with success status, data or error information
 */
export async function loginUser(credentials: {
  login: string;
  password: string;
}): Promise<loginResult> {
  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    const data: loginResponse = await response.json();

    if (!response.ok) {
      return {
        success: false,
        errors: data.errors || {},
      };
    }

    if (data.access_token && data.refresh_token) {
      saveTokensToCookies(data.access_token, data.refresh_token);
    }

    return {
      success: true,
    };
  } catch (error: any) {
    console.error('Login error:', error);
    return {
      success: false,
      errors: {
        other: error.message || 'An unexpected error occurred during registration'
      },
    };
  }
}

/**
 * Logs out the user by clearing tokens
 */
export async function logoutUser() {
  const cookieStore = await cookies();
  cookieStore.delete('access_token');
  cookieStore.delete('refresh_token');

  return {
    success: true,
  };
}
