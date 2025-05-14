'use client';

import {
  syncFiltersWithParams,
  syncRangesWithParams,
} from '@/app/lib/(home)/syncWithParams';
import { FilterFieldData, RangeFieldData } from '@/app/lib/definitions';
import { BaseFilterTemplate } from '@/app/ui/(home)/base-filter-template/base-filter-template';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { BaseRangeTemplate } from './base-filter-template/base-range-template';

export default function Filters({
  filtersData,
  rangesData,
}: {
  filtersData: FilterFieldData[];
  rangesData: RangeFieldData[];
}) {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const filters = syncFiltersWithParams(filtersData, searchParams);
  const ranges = syncRangesWithParams(rangesData, searchParams);

  function handleFilterChange(name: string, selected: string[]) {
    const params = new URLSearchParams(searchParams);
    const sanitizedName = name.replace(/\s+/g, '');

    params.set('page', '1');
    if (selected.length > 0) {
      params.set(sanitizedName, selected.join(','));
    } else {
      params.delete(sanitizedName);
    }

    replace(`${pathname}?${params.toString()}`);
  }

  function handleRangeChange(
    name: string,
    range: { min: number; max: number },
  ) {
    const params = new URLSearchParams(searchParams);
    const sanitizedName = name.replace(/\s+/g, '');

    params.set('page', '1');
    if (range.min !== 0) {
      params.set(`${sanitizedName}_min`, range.min.toString());
    } else {
      params.delete(`${sanitizedName}_min`);
    }

    if (range.max !== 0) {
      params.set(`${sanitizedName}_max`, range.max.toString());
    } else {
      params.delete(`${sanitizedName}_max`);
    }

    replace(`${pathname}?${params.toString()}`);
  }

  return (
    <div>
      <h2 className='mb-2 text-sm font-semibold text-gray-700'>Filters</h2>
      <div className='space-y-3'>
        {filters.map((filter, index) => (
          <BaseFilterTemplate
            key={index}
            name={filter.fieldName}
            options={filter.options}
            selected={filter.selected}
            onChange={handleFilterChange}
          />
        ))}
        {ranges.map((range, index) => (
          <BaseRangeTemplate
            key={index}
            fieldName={range.fieldName}
            rangeBase={range.range}
            onChange={handleRangeChange}
          />
        ))}
      </div>
    </div>
  );
}
