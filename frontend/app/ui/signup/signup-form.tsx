'use client';

import { signup } from '@/app/actions/signup';
import { SignupFormState } from '@/app/lib/definitions/signup';
import { ArrowRightIcon } from '@heroicons/react/20/solid';
import {
  AtSymbolIcon,
  BuildingOfficeIcon,
  IdentificationIcon,
  KeyIcon,
  UserCircleIcon,
  UserIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useActionState, useState } from 'react';
import { Button } from '../(common)/button';
import InputField from '../(common)/input-field';

export default function SignupForm({
  baseAccountType,
}: {
  baseAccountType: string;
}) {
  const initialState: SignupFormState = {
    errors: {},
    values: {},
  };

  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();

  const [backupAccountType, setBackupAccountType] = useState<string | null>(
    null,
  );

  const [accountType, setAccountType] = useState<string>(
    searchParams.get('accountType') || backupAccountType || baseAccountType,
  );

  function changeAccountType(newAccountType: string) {
    const params = new URLSearchParams(searchParams);
    setAccountType(newAccountType);
    params.set('accountType', newAccountType.toString());
    router.replace(`${pathname}?${params.toString()}`);
  }

  const [state, action] = useActionState(signup, initialState);

  const handleSubmit = (formData: FormData) => {
    formData.append('selector', accountType);
    setBackupAccountType(accountType);
    return action(formData);
  };

  return (
    <form className='space-y-3' action={handleSubmit} key={accountType}>
      <div className='flex-1 rounded-lg bg-gray-50 px-6 pt-8 pb-4'>
        <h1 className='mb-3 text-2xl'>Sign up</h1>
        <div className='w-full'>
          <div className='mb-4'>
            <label className='mb-3 block text-xs font-medium text-gray-900'>
              Account Type
            </label>
            <div className='flex gap-4'>
              <div className='flex items-center'>
                <input
                  id='personal'
                  type='radio'
                  value='P'
                  checked={accountType === 'P'}
                  onChange={() => changeAccountType('P')}
                  className='h-4 w-4 text-blue-500 focus:ring-blue-400'
                />
                <label
                  htmlFor='personal'
                  className='ml-2 text-sm font-medium text-gray-900'
                >
                  Personal
                </label>
              </div>
              <div className='flex items-center'>
                <input
                  id='company'
                  type='radio'
                  value='C'
                  checked={accountType === 'C'}
                  onChange={() => changeAccountType('C')}
                  className='h-4 w-4 text-blue-500 focus:ring-blue-400'
                />
                <label
                  htmlFor='business'
                  className='ml-2 text-sm font-medium text-gray-900'
                >
                  Company
                </label>
              </div>
            </div>
            {/* <input type="hidden" name="selector" value={accountType} />  Wartość selector */}
          </div>
          <InputField
            id='username'
            name='username'
            label='Username'
            placeholder='Enter your username'
            defaultValue={state?.values?.username || ''}
            icon={IdentificationIcon}
            errors={state?.errors?.username || []}
            required
          />

          <InputField
            id='email'
            name='email'
            label='Email'
            placeholder='Enter your email address'
            defaultValue={state?.values?.email || ''}
            icon={AtSymbolIcon}
            errors={state?.errors?.email || []}
            required
          />

          {accountType === 'P' && (
            <>
              <InputField
                id='person_name'
                name='person_name'
                label='Name'
                placeholder='Enter your name'
                defaultValue={state?.values?.person_name || ''}
                icon={UserIcon}
                errors={state?.errors?.person_name || []}
                required
              />
              <InputField
                id='person_surname'
                name='person_surname'
                label='Surname'
                placeholder='Enter your surname'
                defaultValue={state?.values?.person_surname || ''}
                icon={UserCircleIcon}
                errors={state?.errors?.person_surname || []}
                required
              />
            </>
          )}

          {accountType === 'C' && (
            <>
              <InputField
                id='company_name'
                name='company_name'
                label='Company Name'
                placeholder='Enter your company name'
                defaultValue={state?.values?.company_name || ''}
                icon={BuildingOfficeIcon}
                errors={state?.errors?.company_name || []}
                required
              />
              <InputField
                id='company_nip'
                name='company_nip'
                label='NIP (Tax ID)'
                placeholder='Enter your NIP'
                defaultValue={state?.values?.company_nip || ''}
                icon={BuildingOfficeIcon}
                errors={state?.errors?.company_nip || []}
                required
              />
            </>
          )}

          <InputField
            id='password'
            name='password'
            label='Password'
            placeholder='Enter password'
            type='password'
            defaultValue=''
            icon={KeyIcon}
            errors={state?.errors?.password || []}
            required
          />

          <InputField
            id='confirm_password'
            name='confirm_password'
            label='Confirm Password'
            placeholder='Confirm your password'
            type='password'
            defaultValue=''
            icon={KeyIcon}
            errors={state?.errors?.confirm_password || []}
            required
          />
        </div>

        <Button type='submit' className='mt-6 w-full'>
          Sign up <ArrowRightIcon className='ml-auto h-5 w-5 text-gray-50' />
        </Button>

        <Link href='/login' className='mt-2 block text-xs font-medium'>
          Already have an account? Log in
        </Link>
      </div>
    </form>
  );
}
