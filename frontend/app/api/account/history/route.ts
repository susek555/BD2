import { HistoryOffer } from '@/app/lib/definitions/SaleOffer';
import { NextResponse } from 'next/server';

export async function GET(request: Request) {
  const offers: HistoryOffer[] = [
    {
      id: '1',
      name: 'Audi A4',
      productionYear: 2000,
      mileage: 150000,
      color: 'Green',
      price: 10000,
      isAuction: true,
      dateEnd: '2025-05-24',
    },
    {
      id: '2',
      name: 'Volkswagen Golf',
      productionYear: 2005,
      mileage: 120000,
      color: 'Blue',
      price: 15000,
      isAuction: false,
      dateEnd: '2025-04-11',
    },
    {
      id: '3',
      name: 'Porsche 911',
      productionYear: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      isAuction: true,
      dateEnd: '2024-11-11',
    },
    {
      id: '4',
      name: 'Porsche 911',
      productionYear: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      isAuction: true,
      dateEnd: '2024-11-11',
    },
    {
      id: '5',
      name: 'Porsche 911',
      productionYear: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      isAuction: true,
      dateEnd: '2024-11-11',
    },
    {
      id: '6',
      name: 'Porsche 911',
      productionYear: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      isAuction: true,
      dateEnd: '2024-11-11',
    },
  ];
  const pagination = {
    total_pages: 1,
    total_records: 6,
  };

  return NextResponse.json({
    offers,
    pagination
  });
}
