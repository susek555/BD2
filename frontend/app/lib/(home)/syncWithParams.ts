import { FilterFieldData, RangeFieldData } from "@/app/lib/definitions";

export function syncFiltersWithParams(
    filters: FilterFieldData[],
    searchParams: URLSearchParams,
  ): FilterFieldData[] {
    return filters.map((filter) => {
      const selectedValues = searchParams.get(filter.fieldName)?.split(",") ?? [];
      return { ...filter, selected: selectedValues };
    });
  }

  // Synchronizacja zakresów z URL
  export function syncRangesWithParams(
    ranges: RangeFieldData[],
    searchParams: URLSearchParams,
  ): RangeFieldData[] {
    return ranges.map((range) => {
      const min = parseInt(searchParams.get(`${range.fieldName}_min`) ?? "0", 10);
      const max = parseInt(searchParams.get(`${range.fieldName}_max`) ?? "0", 10);
      return { ...range, range: { min, max } };
    });
  }

