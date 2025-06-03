// SearchParams

export type SearchParams = {
  query?: string;
  pagination: {
    page: number;
    page_size: number;
  }
  order_key?: string;
  is_order_desc?: boolean;
  offer_type?: string;
  manufacturers?: string[];
  models?: string[];
  drives?: string[];
  colors?: string[];
  transmissions?: string[];
  fuel_types?: string[];
  price_range?: {
    min?: number;
    max?: number;
  };
  mileage_range?: {
    min?: number;
    max?: number;
  };
  year_range?: {
    min?: number;
    max?: number;
  };
  engine_power_range?: {
    min?: number;
    max?: number;
  };
  engine_capacity_range?: {
    min?: number;
    max?: number;
  };
  // offer_creation_date_range: {
  //   min: string | null;
  //   max: string | null;
  // };
}

export const parseIntOrUndefined = (
  value: string | undefined,
): number | undefined => {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  return isNaN(parsed) ? undefined : parsed;
};

export const parseArrayOrUndefined = (
  value: string | string[] | undefined,
): string[] | undefined => {
  if (!value) return undefined;
  if (Array.isArray(value)) {
    return value.flatMap((v) => v.split(',').map(s => s.trim()).filter(Boolean));
  }
  return value.split(',').map(s => s.trim()).filter(Boolean);
};

export const trimAllAfterFirstSpace = (
  values: string[] | undefined,
): string[] | undefined => {
  if (!values) return undefined;
  return values.map(value => {
    const spaceIndex = value.indexOf(' ');
    if (spaceIndex === -1) return value;
    return value.substring(spaceIndex + 1);
  });
};


export function parseFiltersParams(searchParams: {
  query?: string;
  page?: string;
  sortKey?: string;
  isSortDesc?: string;
  offerType?: string;
  Producers?: string[];
  Models?: string[];
  Colors?: string[];
  Drivetypes?: string[];
  Gearboxes?: string[];
  Fueltypes?: string[];
  Price_min?: string;
  Price_max?: string;
  Mileage_min?: string;
  Mileage_max?: string;
  Productionyear_min?: string;
  Productionyear_max?: string;
  Enginecapacity_min?: string;
  Enginecapacity_max?: string;
  Enginepower_min?: string;
  Enginepower_max?: string;
} | undefined) : SearchParams {
  return {
      query: searchParams?.query,
      pagination: {
        page: searchParams?.page ? parseInt(searchParams.page, 10) : 1,
        page_size: 6,
      },
      order_key: searchParams?.sortKey,
      is_order_desc: searchParams?.isSortDesc === "true" ? true : false,
      offer_type: searchParams?.offerType,
      manufacturers: parseArrayOrUndefined(searchParams?.Producers),
      models: trimAllAfterFirstSpace(parseArrayOrUndefined(searchParams?.Models)),
      colors: parseArrayOrUndefined(searchParams?.Colors),
      transmissions: parseArrayOrUndefined(searchParams?.Gearboxes),
      fuel_types: parseArrayOrUndefined(searchParams?.Fueltypes),
      drives: parseArrayOrUndefined(searchParams?.Drivetypes),
      price_range: parseIntOrUndefined(searchParams?.Price_min) || parseIntOrUndefined(searchParams?.Price_max) ? {
        min: parseIntOrUndefined(searchParams?.Price_min),
        max: parseIntOrUndefined(searchParams?.Price_max),
      } : undefined,
      mileage_range: parseIntOrUndefined(searchParams?.Mileage_min) || parseIntOrUndefined(searchParams?.Mileage_max) ? {
        min: parseIntOrUndefined(searchParams?.Mileage_min),
        max: parseIntOrUndefined(searchParams?.Mileage_max),
      } : undefined,
      year_range: parseIntOrUndefined(searchParams?.Productionyear_min) || parseIntOrUndefined(searchParams?.Productionyear_max) ? {
        min: parseIntOrUndefined(searchParams?.Productionyear_min),
        max: parseIntOrUndefined(searchParams?.Productionyear_max),
      } : undefined,
      engine_capacity_range: parseIntOrUndefined(searchParams?.Enginecapacity_min) || parseIntOrUndefined(searchParams?.Enginecapacity_max) ? {
        min: parseIntOrUndefined(searchParams?.Enginecapacity_min),
        max: parseIntOrUndefined(searchParams?.Enginecapacity_max),
      } : undefined,
      engine_power_range: parseIntOrUndefined(searchParams?.Enginepower_min) || parseIntOrUndefined(searchParams?.Enginepower_max) ? {
        min: parseIntOrUndefined(searchParams?.Enginepower_min),
        max: parseIntOrUndefined(searchParams?.Enginepower_max),
      } : undefined,
    };
}
