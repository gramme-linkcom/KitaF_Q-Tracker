// 型定義
export interface QueueStatus {
  waitTime: number;
  timeRequired: number;
  waitingGroups: number;
  myAheadGroups : number;
  currentNumber: number;
  nextNumber: number;
  isActive: boolean;
  isBookingAvailable: boolean;
  noticeMessage: string;
  infoMessage?: string;
}

// 現在の状況を取得する (GET)
export async function fetchQueueStatus(bookingNumber: number): Promise<QueueStatus> {
  let url = "/api/data";
  if (bookingNumber > 0) {
    url += `?myNumber=${bookingNumber}`;
  }

  const response = await fetch(url);
  
  if (!response.ok) throw new Error("データの取得に失敗しました");
  return response.json();
}

// 整理券を発行する (POST)
export async function bookTicket(pushToken: string = ""): Promise<{ bookingNumber: number; success: boolean }> {
  const response = await fetch("/api/booking", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ pushToken }),
  });
  if (!response.ok) throw new Error("予約に失敗しました");
  return response.json();
}

// 整理券をキャンセルする (POST)
export async function cancelTicket(bookingNumber: number): Promise<{ success: boolean }> {
  const response = await fetch("/api/booking/cancel", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ bookingNumber }),
  });
  if (!response.ok) throw new Error("キャンセルに失敗しました");
  return response.json();
}
