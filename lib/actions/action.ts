export interface Action {
  kind: string;
  options: any;
  execute(message: any[]): Promise<any[]>;
}
