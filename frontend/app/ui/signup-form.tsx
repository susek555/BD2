"use client";

import { ArrowRightIcon } from '@heroicons/react/20/solid';
import {
  AtSymbolIcon,
  BuildingOfficeIcon,
  KeyIcon,
  UserCircleIcon,
  UserIcon
} from '@heroicons/react/24/outline';
import { useState } from 'react';
import { Button } from './button';

export default function SignupForm() {
  const [accountType, setAccountType] = useState('personal');

  return (
    <form className="space-y-3">
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
                  name="accountType"
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
                  name="accountType"
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
                required
              />
              <AtSymbolIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
          </div>

          {accountType === 'personal' && (
            <>
              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="name"
                >
                  Name
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="name"
                    type="text"
                    name="name"
                    placeholder="Enter your name"
                    required
                  />
                  <UserIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
              </div>

              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="surname"
                >
                  Surname
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="surname"
                    type="text"
                    name="surname"
                    placeholder="Enter your surname"
                    required
                  />
                  <UserCircleIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
              </div>
            </>
          )}

          {accountType === 'business' && (
            <>
              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="businessName"
                >
                  Business Name
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="businessName"
                    type="text"
                    name="businessName"
                    placeholder="Enter your business name"
                    required
                  />
                  <BuildingOfficeIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                </div>
              </div>

              <div className="mb-4">
                <label
                  className="mb-3 block text-xs font-medium text-gray-900"
                  htmlFor="nip"
                >
                  NIP (Tax ID)
                </label>
                <div className="relative">
                  <input
                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                    id="nip"
                    type="text"
                    name="nip"
                    placeholder="Enter your NIP"
                    required
                  />
                  <BuildingOfficeIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
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
                minLength={6}
              />
              <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
          </div>

          <div className="mb-4">
            <label
              className="mb-3 block text-xs font-medium text-gray-900"
              htmlFor="confirmPassword"
            >
              Confirm Password
            </label>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                id="confirmPassword"
                type="password"
                name="confirmPassword"
                placeholder="Confirm your password"
                required
                minLength={6}
              />
              <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
          </div>
        </div>

        <Button className="mt-6 w-full">
          Sign up <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
        </Button>
      </div>
    </form>
  );
}
