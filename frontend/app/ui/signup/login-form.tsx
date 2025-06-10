'use client';

import { AtSymbolIcon, KeyIcon } from '@heroicons/react/24/outline';
import { ArrowRightIcon } from '@heroicons/react/24/solid';
import { signIn } from 'next-auth/react';
import Link from 'next/link';
import { redirect, useSearchParams } from 'next/navigation';
import { useState } from 'react';
import { Button } from '../(common)/button';
import InputField from '../(common)/input-field';
import { LoginError } from '../../lib/definitions/signup';

export default function LoginForm({
  onLoginSuccess,
}: {
  onLoginSuccess?: () => void;
}) {
  const searchParams = useSearchParams();
  const [errors, setErrors] = useState<LoginError>({});

  const callbackUrl = searchParams.get('callbackUrl') || '/';

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setErrors({});

    const formData = new FormData(e.currentTarget);
    const login = formData.get('login') as string;
    const password = formData.get('password') as string;

    const result = await signIn('credentials', {
      redirect: false,
      login,
      password,
      callbackUrl,
    });

    console.log('Login result:', result);

    if (result?.error) {
      try {
        setErrors(JSON.parse(result.error));
      } catch {
        setErrors({ server: [result.error] });
      }
      console.log('Parsed login errors:', errors);
    } else {
      if (onLoginSuccess) {
        onLoginSuccess();
      }

      redirect(callbackUrl);
    }
  }

  return (
    <form className='space-y-3' onSubmit={handleSubmit}>
      <div className='flex-1 rounded-lg bg-gray-50 px-6 pt-8 pb-4'>
        <h1 className='mb-3 text-2xl'>Log in</h1>
        <div className='w-full'>
          <InputField
            id='login'
            name='login'
            label='Login'
            placeholder='Enter your email address'
            defaultValue=''
            icon={AtSymbolIcon}
            errors={[]}
            required
          />

          <InputField
            id='password'
            name='password'
            label='Password'
            placeholder='Enter password'
            type='password'
            defaultValue=''
            icon={KeyIcon}
            errors={[]}
            required
          />
        </div>
        <Button className='mt-4 w-full' type='submit'>
          Log in <ArrowRightIcon className='ml-auto h-5 w-5 text-gray-50' />
        </Button>
        <Link href='/signup' className='mt-2 block text-xs font-medium'>
          Need an account? Sign up
        </Link>
        <div className='mt-4 flex h-8 items-end space-x-1'>
          <div id='credentials-error' aria-live='polite' aria-atomic='true'>
            {errors.credentials &&
              errors.credentials.map((error: string) => (
                <p className='mt-2 text-sm text-red-500' key={error}>
                  {error}
                </p>
              ))}
          </div>
          <div id='server-error' aria-live='polite' aria-atomic='true'>
            {errors.server &&
              errors.server.map((error: string) => (
                <p className='mt-2 text-sm text-red-500' key={error}>
                  {error}
                </p>
              ))}
          </div>
        </div>
      </div>
    </form>
  );
}
