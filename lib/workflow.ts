import { Action } from "./actions/action";
import { Trigger } from "./triggers/trigger";

export interface FlowNode {
  action: Action;
  nodes: FlowNode[];
}

export interface Flow {
  trigger: Trigger;
  root: FlowNode;
}
