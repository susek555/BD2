"use server";


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

    console.log("Registration response: ", data);

    if (!response.ok) {
      return {
        success: false,
        errors: data.errors || {},
      };
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
