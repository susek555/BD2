'use client';

import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';

export default function Sorting({
  sortingOptions,
}: {
  sortingOptions: string[];
}) {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const [isKeyDropdownOpen, setIsKeyDropdownOpen] = useState(false);
  const [isOrderDropdownOpen, setIsOrderDropdownOpen] = useState(false);

  const selectedKey = searchParams.get('sortKey') || 'Base';
  const isDescSelected = searchParams.get('isSortDesc') || false;

  const toggleKeyDropdown = () => setIsKeyDropdownOpen(!isKeyDropdownOpen);
  const toggleOrderDropdown = () =>
    setIsOrderDropdownOpen(!isOrderDropdownOpen);

  const handleKeyChange = (key: string) => {
    const params = new URLSearchParams(searchParams);
    if (key === 'Base') {
      params.delete('sortKey');
    } else {
      params.set('sortKey', key);
    }
    params.set('page', '1');
    replace(`${pathname}?${params.toString()}`);
    setIsKeyDropdownOpen(false);
  };

  const handleOrderChange = (order: string) => {
    const params = new URLSearchParams(searchParams);
    if (order === 'Desc') {
      params.set('isSortDesc', 'true');
    } else {
      params.delete('isSortDesc');
    }
    params.set('page', '1');
    replace(`${pathname}?${params.toString()}`);
    setIsOrderDropdownOpen(false);
  };

  return (
    <div>
      <h2 className='mb-2 text-sm font-semibold text-gray-700'>Order by</h2>
      <div className='flex gap-2'>
        <div className='relative flex-[2]'>
          <button
            className='flex w-full items-center justify-between rounded-lg border bg-white px-4 py-2 transition-colors hover:bg-gray-50'
            onClick={toggleKeyDropdown}
          >
            <span className='text-gray-900'>{selectedKey}</span>
            <span className='text-sm text-gray-500'>
              {isKeyDropdownOpen ? '▲' : '▼'}
            </span>
          </button>
          {isKeyDropdownOpen && (
            <div className='absolute top-full right-0 left-0 z-10 mt-1 rounded-lg border bg-white shadow-lg'>
              {sortingOptions.map((key) => (
                <div
                  key={key}
                  className='cursor-pointer px-4 py-2 first:rounded-t-lg last:rounded-b-lg hover:bg-gray-100'
                  onClick={() => handleKeyChange(key)}
                >
                  <span className='text-gray-900'>{key}</span>
                </div>
              ))}
            </div>
          )}
        </div>
        <div className='relative flex-1'>
          <button
            className='flex w-full items-center justify-between rounded-lg border bg-white px-4 py-2 transition-colors hover:bg-gray-50'
            onClick={toggleOrderDropdown}
          >
            <span className='text-gray-900'>
              {isDescSelected === 'true' ? 'Desc' : 'Asc'}
            </span>
            <span className='text-sm text-gray-500'>
              {isOrderDropdownOpen ? '▲' : '▼'}
            </span>
          </button>
          {isOrderDropdownOpen && (
            <div className='absolute top-full right-0 left-0 z-10 mt-1 rounded-lg border bg-white shadow-lg'>
              {['Asc', 'Desc'].map((order) => (
                <div
                  key={order}
                  className='cursor-pointer px-4 py-2 first:rounded-t-lg last:rounded-b-lg hover:bg-gray-100'
                  onClick={() => handleOrderChange(order)}
                >
                  <span className='text-gray-900'>{order}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
