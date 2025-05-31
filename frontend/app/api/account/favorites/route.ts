import { SaleOffer } from '@/app/lib/definitions/SaleOffer';
import { NextResponse } from 'next/server';

export async function GET(request: Request) {
  const offers: SaleOffer[] = [
    {
      id: '1',
      name: 'Audi A4',
      production_year: 2000,
      mileage: 150000,
      color: 'Green',
      price: 10000,
      is_auction: true,
      is_liked: true,
    },
    {
      id: '2',
      name: 'Volkswagen Golf',
      production_year: 2005,
      mileage: 120000,
      color: 'Blue',
      price: 15000,
      is_auction: false,
      is_liked: true,
    },
    {
      id: '3',
      name: 'Porsche 911',
      production_year: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      is_auction: true,
      is_liked: true,
    },
  ];
  const pagination = {
    total_pages: 1,
    total_records: 3,
  };

  return NextResponse.json({
    offers,
    pagination
  });
}
