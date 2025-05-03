"use client";

import { signup } from '@/app/actions/auth';
import { ArrowRightIcon } from '@heroicons/react/20/solid';
import {
  AtSymbolIcon,
  BuildingOfficeIcon,
  IdentificationIcon,
  KeyIcon,
  UserCircleIcon,
  UserIcon
} from '@heroicons/react/24/outline';
import { useActionState, useState } from 'react';
import { SignupFormState } from '../lib/definitions';
import { Button } from './button';

const initialState: SignupFormState = {
  errors: {},
  values: {},
};


export default function SignupForm() {
  const initialState: SignupFormState = {
    errors: {},
    values: {},
  };

  const [state, action] = useActionState(signup, initialState);

  const [accountType, setAccountType] = useState(
    state?.values?.selector || 'personal'
  );

  return (
    <form className="space-y-3" action={action}>
      <div className="flex-1 rounded-lg bg-gray-50 px-6 pb-4 pt-8">
        <h1 className="mb-3 text-2xl">
          Sign up
        </h1>
        <div className="w-full">
          <div className="mb-4">
            <label className="mb-3 block text-xs font-medium text-gray-900">
              Account Type
            </label>
            <div className="flex gap-4">
              <div className="flex items-center">
                <input
                  id="personal"
                  type="radio"
                  name="selector"
                  value="personal"
                  checked={accountType === 'personal'}
                  onChange={() => setAccountType('personal')}
                  className="h-4 w-4 text-blue-500 focus:ring-blue-400"
                />
                <label htmlFor="personal" className="ml-2 text-sm font-medium text-gray-900">
                  Personal
                </label>
              </div>
              <div className="flex items-center">
                <input
                  id="business"
                  type="radio"
                  name="selector"
                  value="business"
                  checked={accountType === 'business'}
                  onChange={() => setAccountType('business')}
                  className="h-4 w-4 text-blue-500 focus:ring-blue-400"
                />
                <label htmlFor="business" className="ml-2 text-sm font-medium text-gray-900">
                  Business
                </label>
              </div>
            </div>
          </div>

          <div className="mb-4">
            <label
              className="mb-3 block text-xs font-medium text-gray-900"
              htmlFor="username"
            >
              Username
            </label>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                id="username"
                type="text"
                name="username"
                placeholder="Enter your username"
                defaultValue={state?.values?.username}
                required
              />
              <IdentificationIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
            <div id="username-error" aria-live="polite" aria-atomic="true">
              {state?.errors?.username &&
                state.errors.username.map((error: string) => (
                  <p className="mt-2 text-sm text-red-500" key={error}>
                    {error}
                  </p>
                ))}
            </div>
          </div>

          <div className="mb-4">
            <label
              className="mb-3 block text-xs font-medium text-gray-900"
              htmlFor="email"
            >
              Email
            </label>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                id="email"
                type="text"
                name="email"
                placeholder="Enter your email address"
                defaultValue={state?.values?.email}
                required
              />
              <AtSymbolIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
            <div id="email-error" aria-live="polite" aria-atomic="true">
              {state?.errors?.email &&
                state.errors.email.map((error: string) => (
                  <p className="mt-2 text-sm text-red-500" key={error}>
                    {error}
                  </p>
                ))}
            </div>
          </div>

          {accountType === 'personal' && (
            <>
              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="person_name"
                >
                  Name
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="person_name"
                    type="text"
                    name="person_name"
                    placeholder="Enter your name"
                    defaultValue={state?.values?.person_name}
                    required
                  />
                  <UserIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
                <div id="person-name-error" aria-live="polite" aria-atomic="true">
                  {state?.errors?.person_name &&
                    state.errors.person_name.map((error: string) => (
                      <p className="mt-2 text-sm text-red-500" key={error}>
                        {error}
                      </p>
                    ))}
                </div>
              </div>

              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="person_surname"
                >
                  Surname
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="person_surname"
                    type="text"
                    name="person_surname"
                    placeholder="Enter your surname"
                    defaultValue={state?.values?.person_surname}
                    required
                  />
                  <UserCircleIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
                <div id="person-surname-error" aria-live="polite" aria-atomic="true">
                  {state?.errors?.person_surname &&
                    state.errors.person_surname.map((error: string) => (
                      <p className="mt-2 text-sm text-red-500" key={error}>
                        {error}
                      </p>
                    ))}
                </div>
              </div>
            </>
          )}

          {accountType === 'business' && (
            <>
              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="business_name"
                >
                  Business Name
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="business_name"
                    type="text"
                    name="business_name"
                    placeholder="Enter your business name"
                    defaultValue={state?.values?.business_name}
                    required
                  />
                  <BuildingOfficeIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
                <div id="business-name-error" aria-live="polite" aria-atomic="true">
                  {state?.errors?.business_name &&
                    state.errors.business_name.map((error: string) => (
                      <p className="mt-2 text-sm text-red-500" key={error}>
                        {error}
                      </p>
                    ))}
                </div>
              </div>

              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="business_nip"
                >
                  NIP (Tax ID)
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="business_nip"
                    type="text"
                    name="business_nip"
                    placeholder="Enter your NIP"
                    defaultValue={state?.values?.business_nip}
                    required
                  />
                  <BuildingOfficeIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
                <div id="business-nip-error" aria-live="polite" aria-atomic="true">
                  {state?.errors?.business_nip &&
                    state.errors.business_nip.map((error: string) => (
                      <p className="mt-2 text-sm text-red-500" key={error}>
                        {error}
                      </p>
                    ))}
                </div>
              </div>
            </>
          )}

          <div className="mb-4">
            <label
              className="mb-3 block text-xs font-medium text-gray-900"
              htmlFor="password"
            >
              Password
            </label>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                id="password"
                type="password"
                name="password"
                placeholder="Enter password"
                required
              />
              <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
            <div id="password-error" aria-live="polite" aria-atomic="true">
              {state?.errors?.password && (
                <div>
                  <p className="mt-2 text-sm text-red-500">Password must:</p>
                  <ul>
                    {state.errors.password.map((error) => (
                      <li className="text-sm text-red-500" key={error}>- {error}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>

          <div className="mb-4">
            <label
              className="mb-3 block text-xs font-medium text-gray-900"
              htmlFor="confirm_password"
            >
              Confirm Password
            </label>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                id="confirm_password"
                type="password"
                name="confirm_password"
                placeholder="Confirm your password"
                required
              />
              <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
            <div id="confirm-password-error" aria-live="polite" aria-atomic="true">
              {state?.errors?.confirm_password &&
                state.errors.confirm_password.map((error: string) => (
                  <p className="mt-2 text-sm text-red-500" key={error}>
                    {error}
                  </p>
                ))}
            </div>
          </div>
        </div>

        <Button type="submit" className="mt-6 w-full">
          Sign up <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
        </Button>
      </div>
    </form>
  );
}
