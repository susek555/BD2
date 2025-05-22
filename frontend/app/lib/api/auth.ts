'use server';

import { API_URL } from '../constants';

export type registerResult = {
  errors?: {
    username?: string[];
    email?: string[];
    nip?: string[];
    other?: string[];
  };
};

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

    if (response.ok && response.headers.get('Content-Length') === '0') {
      return {};
    }

    const data: registerResult = await response.json();

    console.log('Registration response: ', data);

    if (!response.ok) {
      return {
        errors: data.errors || {},
      };
    }

    return {};
  } catch (error) {
    console.error('Registration error:', error);
    return {
      errors: {
        other: [error as string],
      },
    };
  }
}
