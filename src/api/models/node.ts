import { Stream } from './stream';
import { Team } from './team';

export enum NodeGroupType {
    Invalid = "INVALID",
    Organizational = "ORGANIZATIONAL",
    RoundRobin = "ROUND_ROBIN",
    Swiss = "SWISS",
    BracketSingle = "BRACKET_SINGLE",
    BracketDoubleSeedLoser = "BRACKET_DOUBLE_SEED_LOSER",
    BracketDoubleAllWinner = "BRACKET_DOUBLE_ALL_WINNER",
    Showmatch = "SHOWMATCH",
    GSL = "GSL"
};

export enum NodeType {
    Invalid = "INVALID",
    BestOfOne = "BEST_OF_ONE",
    BestOfThree = "BEST_OF_THREE",
    BestOfFive = "BEST_OF_FIVE",
    BestOfTwo = "BEST_OF_TWO"
};

export type Node = {
    id: number;
    scheduledTime?: number;
    actualTime?: number;
    nodeType?: NodeType;
    hasStarted?: boolean;
    isCompleted?: boolean;
    teamOne?: Team;
    teamTwo?: Team;
    streams?: Stream[];
};

export type NodeGroup = {
    id: number;
    name?: string;
    nodeGroupType?: NodeGroupType;
    round?: number;
    nodes?: Node[];
};