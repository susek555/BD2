'use client';

import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';

export default function OfferType() {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const options = ['All', 'Auctions', 'Buy Now'];
  const selectedOption = searchParams.get('offerType') || 'All';

  const toggleDropdown = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  const handleOptionChange = (option: string) => {
    const params = new URLSearchParams(searchParams);
    if (option === 'All') {
      params.delete('offerType');
    } else {
      params.set('offerType', option);
    }
    params.set('page', '1');
    replace(`${pathname}?${params.toString()}`);
    setIsDropdownOpen(false);
  };

  return (
    <div>
      <h2 className='mb-2 text-sm font-semibold text-gray-700'>Offer Type</h2>
      <div className='relative'>
        <button
          className='flex w-full items-center justify-between rounded-lg border bg-white px-4 py-2 transition-colors hover:bg-gray-50'
          onClick={toggleDropdown}
        >
          <span className='text-gray-900'>{selectedOption}</span>
          <span className='text-sm text-gray-500'>
            {isDropdownOpen ? '▲' : '▼'}
          </span>
        </button>
        {isDropdownOpen && (
          <div className='absolute top-full right-0 left-0 z-10 mt-1 rounded-lg border bg-white shadow-lg'>
            {options.map((option) => (
              <div
                key={option}
                className='cursor-pointer px-4 py-2 first:rounded-t-lg last:rounded-b-lg hover:bg-gray-100'
                onClick={() => handleOptionChange(option)}
              >
                <span className='text-gray-900'>{option}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
